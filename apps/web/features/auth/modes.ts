/** Text for the two auth pages. One map, so a mode is added in one place. */
export const authModes = {
  "sign-in": {
    title: "Sign in",
    description: "Contribute, review and manage your API keys.",
    submit: "Sign in",
    switchText: "New here?",
    switchLink: "Create an account",
    switchHref: "/sign-up",
  },
  "sign-up": {
    title: "Create an account",
    description: "Free for developers, contributors and reviewers.",
    submit: "Create account",
    switchText: "Already have an account?",
    switchLink: "Sign in",
    switchHref: "/sign-in",
  },
} as const

export type AuthMode = keyof typeof authModes
