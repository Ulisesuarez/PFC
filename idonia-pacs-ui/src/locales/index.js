export const locales = [
  {
    code: 'es-ES',
    shortCode: 'es',
    name: 'Español'
  },
  {
    code: 'ca-ES',
    shortCode: 'ca',
    name: 'Català'
  },
  {
    code: 'en-US',
    shortCode: 'en',
    name: 'English'
  }
]
export const defaultLocale = locales.find(({ code }) => code === 'en-US')
export const parseLocale = someString => {
  const normalize = value => value.toLowerCase().replace('_', '').replace('_', '')
  const normalized = normalize(someString)
  return locales.find(({ code }) => normalize(code).startsWith(normalized))
}

export default locales
