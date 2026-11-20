import next from "@vernacular/eslint-config/next"

const config = [...next, { ignores: ["components/ui/**"] }]

export default config
