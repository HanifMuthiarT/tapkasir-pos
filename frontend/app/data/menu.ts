export type MenuCategory = "Kopi" | "Non Kopi" | "Teh" | "Camilan";

export type MenuItem = {
  id: string;
  name: string;
  category: MenuCategory;
  description: string;
  price: number;
  image: string;
  available: boolean;
  tag?: string;
};

export const categories = ["Semua Menu", "Kopi", "Non Kopi", "Teh", "Camilan"] as const;
export type CategoryFilter = (typeof categories)[number];

export const menuItems: MenuItem[] = [
  {
    id: "espresso",
    name: "Espresso",
    category: "Kopi",
    description: "Espresso klasik dengan karakter kuat dan aroma yang pekat.",
    price: 18000,
    image: "https://images.unsplash.com/photo-1511920170033-f8396924c348?w=800&q=85",
    available: true,
  },
  {
    id: "iced-americano",
    name: "Iced Americano",
    category: "Kopi",
    description: "Espresso dan air dingin, sederhana dan menyegarkan.",
    price: 24000,
    image: "https://images.unsplash.com/photo-1461023058943-07fcbe16d735?w=800&q=85",
    available: true,
    tag: "Favorit",
  },
  {
    id: "cappuccino",
    name: "Cappuccino",
    category: "Kopi",
    description: "Perpaduan espresso, susu hangat, dan foam lembut.",
    price: 28000,
    image: "https://images.unsplash.com/photo-1570968915860-54d5c301fa9f?w=800&q=85",
    available: true,
  },
  {
    id: "latte",
    name: "Caffè Latte",
    category: "Kopi",
    description: "Rasa kopi yang halus dengan susu creamy.",
    price: 28000,
    image: "https://images.unsplash.com/photo-1461023058943-07fcbe16d735?w=800&q=85",
    available: true,
  },
  {
    id: "aren",
    name: "Kopi Susu Aren",
    category: "Kopi",
    description: "Kopi susu dengan manis alami gula aren.",
    price: 26000,
    image: "https://images.unsplash.com/photo-1497515114629-f71d768fd07c?w=800&q=85",
    available: true,
    tag: "Best seller",
  },
  {
    id: "mocha",
    name: "Iced Mocha",
    category: "Kopi",
    description: "Kombinasi kopi, susu, dan cokelat yang seimbang.",
    price: 30000,
    image: "https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=800&q=85",
    available: true,
  },
  {
    id: "matcha",
    name: "Matcha Latte",
    category: "Non Kopi",
    description: "Matcha dengan susu lembut dan rasa earthy.",
    price: 29000,
    image: "https://images.unsplash.com/photo-1515823064-d6e0c04616a7?w=800&q=85",
    available: true,
  },
  {
    id: "chocolate",
    name: "Chocolate Milk",
    category: "Non Kopi",
    description: "Cokelat lembut yang nikmat disajikan dingin.",
    price: 27000,
    image: "https://images.unsplash.com/photo-1542990253-0b8be8d5f0ed?w=800&q=85",
    available: true,
  },
  {
    id: "lemon-tea",
    name: "Lemon Tea",
    category: "Teh",
    description: "Teh segar dengan aroma lemon yang ringan.",
    price: 22000,
    image: "https://images.unsplash.com/photo-1544787219-7f47ccb76574?w=800&q=85",
    available: true,
  },
  {
    id: "peach-tea",
    name: "Peach Tea",
    category: "Teh",
    description: "Teh dingin dengan sentuhan buah persik.",
    price: 24000,
    image: "https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=800&q=85",
    available: true,
  },
  {
    id: "croissant",
    name: "Butter Croissant",
    category: "Camilan",
    description: "Croissant berlapis dengan aroma mentega.",
    price: 25000,
    image: "https://images.unsplash.com/photo-1555507036-ab1f4038808a?w=800&q=85",
    available: true,
  },
  {
    id: "cookies",
    name: "Chocolate Cookies",
    category: "Camilan",
    description: "Cookies renyah dengan potongan cokelat.",
    price: 20000,
    image: "https://images.unsplash.com/photo-1499636136210-6f4ee915583e?w=800&q=85",
    available: true,
  },
];

export function rupiah(amount: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amount);
}
