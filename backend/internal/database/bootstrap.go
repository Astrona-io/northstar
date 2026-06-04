package database

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cmdb-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type CategoryManifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Icon        string `yaml:"icon"`
		Description string `yaml:"description"`
		Order       int    `yaml:"order"`
	} `yaml:"spec"`
}

type SubGroupManifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name           string `yaml:"name"`
		ParentCategory string `yaml:"parentCategory"`
	} `yaml:"metadata"`
	Spec struct {
		Icon        string `yaml:"icon"`
		Description string `yaml:"description"`
		Order       int    `yaml:"order"`
	} `yaml:"spec"`
}

type UserManifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Password string `yaml:"password"`
		Role     string `yaml:"role"`
	} `yaml:"spec"`
}

type GroupManifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Description string `yaml:"description"`
	} `yaml:"spec"`
}

type PermissionManifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Effect      string `yaml:"effect"`
		Description string `yaml:"description"`
	} `yaml:"spec"`
}

// ReconcileSeedingYAML parses modular K8s-style manifests under .northstar-data/bootstrap/ (Phase 1 Seeding)
func ReconcileSeedingYAML(db *gorm.DB) error {
	basePath := "../.northstar-data/bootstrap"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// Fallback for handlers test running inside nested package folders
		basePath = "../../.northstar-data/bootstrap"
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			log.Println("[Bootstrap Engine] Warning: .northstar-data/bootstrap directory not found. Skipping YAML seeding.")
			return nil
		}
	}

	categoriesDir := filepath.Join(basePath, "asset-taxonomy", "categories")
	subgroupsDir := filepath.Join(basePath, "asset-taxonomy", "sub-groups")
	usersDir := filepath.Join(basePath, "iam-security", "users")
	groupsDir := filepath.Join(basePath, "iam-security", "groups")
	permissionsDir := filepath.Join(basePath, "iam-security", "permissions")

	// 1. Reconcile Categories (categories/*.yaml)
	if _, err := os.Stat(categoriesDir); err == nil {
		files, err := ioutil.ReadDir(categoriesDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				filePath := filepath.Join(categoriesDir, f.Name())
				b, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}

				decoder := yaml.NewDecoder(bytes.NewReader(b))
				for {
					var m CategoryManifest
					if err := decoder.Decode(&m); err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("[Bootstrap Engine] Error decoding category file %s: %v\n", f.Name(), err)
						return err
					}

					if m.Kind == "Category" && m.Metadata.Name != "" {
						var category models.Category
						err := db.Where("name = ?", m.Metadata.Name).First(&category).Error
						if err != nil && err == gorm.ErrRecordNotFound {
							category = models.Category{
								Name:        m.Metadata.Name,
								Icon:        m.Spec.Icon,
								Description: m.Spec.Description,
								Order:       m.Spec.Order,
							}
							if err := db.Create(&category).Error; err != nil {
								return err
							}
							log.Printf("[Bootstrap Engine] Created Category: %s\n", category.Name)
						} else {
							category.Icon = m.Spec.Icon
							category.Description = m.Spec.Description
							category.Order = m.Spec.Order
							db.Save(&category)
						}
					}
				}
			}
		}
	}

	// 2. Reconcile SubGroups (sub-groups/*.yaml)
	if _, err := os.Stat(subgroupsDir); err == nil {
		files, err := ioutil.ReadDir(subgroupsDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				filePath := filepath.Join(subgroupsDir, f.Name())
				b, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}

				decoder := yaml.NewDecoder(bytes.NewReader(b))
				for {
					var m SubGroupManifest
					if err := decoder.Decode(&m); err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("[Bootstrap Engine] Error decoding subgroup file %s: %v\n", f.Name(), err)
						return err
					}

					if m.Kind == "SubGroup" && m.Metadata.Name != "" && m.Metadata.ParentCategory != "" {
						var parentCat models.Category
						if err := db.Where("name = ?", m.Metadata.ParentCategory).First(&parentCat).Error; err == nil {
							var sub models.SubGroup
							err := db.Where("category_id = ? AND name = ?", parentCat.ID, m.Metadata.Name).First(&sub).Error
							if err != nil && err == gorm.ErrRecordNotFound {
								sub = models.SubGroup{
									CategoryID:  parentCat.ID,
									Name:        m.Metadata.Name,
									Icon:        m.Spec.Icon,
									Description: m.Spec.Description,
									Order:       m.Spec.Order,
								}
								if err := db.Create(&sub).Error; err != nil {
									return err
								}
								log.Printf("[Bootstrap Engine] Created SubGroup: %s &rarr; %s\n", parentCat.Name, sub.Name)
							} else {
								sub.Icon = m.Spec.Icon
								sub.Description = m.Spec.Description
								sub.Order = m.Spec.Order
								db.Save(&sub)
							}
						} else {
							log.Printf("[Bootstrap Engine] Warning: parent category %s not found for subgroup %s. Skipping.\n", m.Metadata.ParentCategory, m.Metadata.Name)
						}
					}
				}
			}
		}
	}

	// 3. Reconcile Permissions (permissions/*.yaml)
	if _, err := os.Stat(permissionsDir); err == nil {
		files, err := ioutil.ReadDir(permissionsDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				filePath := filepath.Join(permissionsDir, f.Name())
				b, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}

				decoder := yaml.NewDecoder(bytes.NewReader(b))
				for {
					var m PermissionManifest
					if err := decoder.Decode(&m); err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("[Bootstrap Engine] Error decoding permission file %s: %v\n", f.Name(), err)
						return err
					}

					if m.Kind == "Permission" && m.Metadata.Name != "" {
						var perm models.Permission
						err := db.Where("name = ?", m.Metadata.Name).First(&perm).Error
						if err != nil && err == gorm.ErrRecordNotFound {
							perm = models.Permission{
								Name:        m.Metadata.Name,
								Effect:      m.Spec.Effect,
								Description: m.Spec.Description,
							}
							if err := db.Create(&perm).Error; err != nil {
								return err
							}
							log.Printf("[Bootstrap Engine] Created Permission: %s (%s)\n", perm.Name, perm.Effect)
						} else {
							perm.Effect = m.Spec.Effect
							perm.Description = m.Spec.Description
							db.Save(&perm)
						}
					}
				}
			}
		}
	}

	// 4. Reconcile Groups (groups/*.yaml)
	if _, err := os.Stat(groupsDir); err == nil {
		files, err := ioutil.ReadDir(groupsDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				filePath := filepath.Join(groupsDir, f.Name())
				b, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}

				decoder := yaml.NewDecoder(bytes.NewReader(b))
				for {
					var m GroupManifest
					if err := decoder.Decode(&m); err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("[Bootstrap Engine] Error decoding group file %s: %v\n", f.Name(), err)
						return err
					}

					if m.Kind == "Group" && m.Metadata.Name != "" {
						var group models.Group
						err := db.Where("name = ?", m.Metadata.Name).First(&group).Error
						if err != nil && err == gorm.ErrRecordNotFound {
							group = models.Group{
								Name:        m.Metadata.Name,
								Description: m.Spec.Description,
							}
							if err := db.Create(&group).Error; err != nil {
								return err
							}
							log.Printf("[Bootstrap Engine] Created Access Group: %s\n", group.Name)
						} else {
							group.Description = m.Spec.Description
							db.Save(&group)
						}
					}
				}
			}
		}
	}

	// 5. Reconcile Users (users/*.yaml)
	if _, err := os.Stat(usersDir); err == nil {
		files, err := ioutil.ReadDir(usersDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				filePath := filepath.Join(usersDir, f.Name())
				b, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}

				decoder := yaml.NewDecoder(bytes.NewReader(b))
				for {
					var m UserManifest
					if err := decoder.Decode(&m); err != nil {
						if err == io.EOF {
							break
						}
						log.Printf("[Bootstrap Engine] Error decoding user file %s: %v\n", f.Name(), err)
						return err
					}

					if m.Kind == "User" && m.Metadata.Name != "" {
						var user models.User
						err := db.Where("username = ?", m.Metadata.Name).First(&user).Error
						if err != nil && err == gorm.ErrRecordNotFound {
							hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Spec.Password), bcrypt.DefaultCost)
							if err != nil {
								return err
							}
							user = models.User{
								Username: m.Metadata.Name,
								Password: string(hashedPassword),
								Role:     m.Spec.Role,
							}
							if err := db.Create(&user).Error; err != nil {
								return err
							}
							log.Printf("[Bootstrap Engine] Created Declarative User: %s (Role: %s)\n", user.Username, user.Role)
						} else {
							// Update role if changed
							user.Role = m.Spec.Role
							db.Save(&user)
						}
					}
				}
			}
		}
	}

	log.Println("[Bootstrap Engine] All modular K8s-style bootstrapping manifests reconciled successfully.")
	return nil
}
