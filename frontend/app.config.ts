// frontend/app.config.ts (Theme Customization matching Gitea / Forgejo Signature Teal and Slate)
export default defineAppConfig({
  ui: {
    primary: 'teal',
    gray: 'slate',
    button: {
      default: {
        size: 'sm'
      },
      rounded: 'rounded-md' // Gitea / Forgejo crisp rounded-md buttons
    },
    card: {
      rounded: 'rounded-md', // Flat panel styling
      shadow: 'shadow-sm',
      border: 'border border-slate-200 dark:border-slate-800'
    },
    modal: {
      rounded: 'rounded-md'
    },
    badge: {
      rounded: 'rounded-md'
    }
  }
})
