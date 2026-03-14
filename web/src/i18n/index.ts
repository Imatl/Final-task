import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './locales/en.json';
import uk from './locales/uk.json';

const getDefaultLanguage = (): string => {
  const saved = localStorage.getItem('language');
  if (saved) return saved;

  const primaryLang = navigator.language || navigator.languages?.[0] || '';
  const isUkOrRu = primaryLang.startsWith('uk') || primaryLang.startsWith('ru');

  return isUkOrRu ? 'uk' : 'en';
};

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: en },
      uk: { translation: uk },
    },
    lng: getDefaultLanguage(),
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
  });

export const changeLanguage = (lang: string) => {
  localStorage.setItem('language', lang);
  i18n.changeLanguage(lang);
};

export default i18n;
