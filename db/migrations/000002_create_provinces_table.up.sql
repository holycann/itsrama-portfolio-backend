CREATE TABLE public.provinces (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying(100) NOT NULL UNIQUE,
    description character varying(500),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone NULL,
    CONSTRAINT provinces_pkey PRIMARY KEY (id)
);

-- Complete list of all 38 provinces in Indonesia (as of 2024)
INSERT INTO public.provinces (id, name, description) VALUES 
  (gen_random_uuid(), 'Aceh', 'Special autonomous province at the northern tip of Sumatra, known for its Islamic culture and tsunami history'),
  (gen_random_uuid(), 'Sumatera Utara', 'North Sumatra province with Lake Toba and diverse ethnic groups including Batak'),
  (gen_random_uuid(), 'Sumatera Selatan', 'South Sumatra province known for Palembang city and Musi River'),
  (gen_random_uuid(), 'Sumatera Barat', 'West Sumatra province, homeland of Minangkabau people and Padang cuisine'),
  (gen_random_uuid(), 'Bengkulu', 'Province on the west coast of Sumatra, birthplace of President Soekarno'),
  (gen_random_uuid(), 'Riau', 'Province rich in oil and gas resources, with Pekanbaru as capital'),
  (gen_random_uuid(), 'Kepulauan Riau', 'Riau Islands province including Batam, strategic location for trade and industry'),
  (gen_random_uuid(), 'Jambi', 'Province in central Sumatra known for rubber and palm oil plantations'),
  (gen_random_uuid(), 'Lampung', 'Southernmost province of Sumatra, gateway between Sumatra and Java'),
  (gen_random_uuid(), 'Kepulauan Bangka Belitung', 'Island province known for tin mining and beautiful beaches'),
  (gen_random_uuid(), 'DKI Jakarta', 'Special Capital Region, Indonesia''s capital and largest metropolitan area'),
  (gen_random_uuid(), 'Banten', 'Westernmost province of Java, includes Serang and industrial areas'),
  (gen_random_uuid(), 'Jawa Barat', 'West Java province, most populous province with Bandung as capital'),
  (gen_random_uuid(), 'Jawa Tengah', 'Central Java province, heartland of Javanese culture with Semarang as capital'),
  (gen_random_uuid(), 'DI Yogyakarta', 'Special Region of Yogyakarta, cultural center and royal sultanate'),
  (gen_random_uuid(), 'Jawa Timur', 'East Java province, second most populous with Surabaya as capital'),
  (gen_random_uuid(), 'Bali', 'Island province famous for Hindu culture, tourism, and beautiful landscapes'),
  (gen_random_uuid(), 'Nusa Tenggara Barat', 'West Nusa Tenggara with Lombok island and Mount Rinjani'),
  (gen_random_uuid(), 'Nusa Tenggara Timur', 'East Nusa Tenggara with Flores, Timor, and Komodo National Park'),
  (gen_random_uuid(), 'Kalimantan Barat', 'West Kalimantan (Borneo) with Pontianak as capital'),
  (gen_random_uuid(), 'Kalimantan Tengah', 'Central Kalimantan with vast forests and Palangka Raya as capital'),
  (gen_random_uuid(), 'Kalimantan Selatan', 'South Kalimantan known for coal mining and floating markets'),
  (gen_random_uuid(), 'Kalimantan Timur', 'East Kalimantan, oil-rich province with Samarinda as capital'),
  (gen_random_uuid(), 'Kalimantan Utara', 'North Kalimantan, newest province established in 2012 with Tanjung Selor as capital'),
  (gen_random_uuid(), 'Sulawesi Utara', 'North Sulawesi with Manado as capital, known for Bunaken marine park'),
  (gen_random_uuid(), 'Sulawesi Tengah', 'Central Sulawesi with diverse ethnic groups and Palu as capital'),
  (gen_random_uuid(), 'Sulawesi Selatan', 'South Sulawesi, homeland of Bugis and Makassar people'),
  (gen_random_uuid(), 'Sulawesi Tenggara', 'Southeast Sulawesi with Kendari as capital and nickel mining'),
  (gen_random_uuid(), 'Gorontalo', 'Small province in northern Sulawesi with unique local culture'),
  (gen_random_uuid(), 'Sulawesi Barat', 'West Sulawesi, separated from South Sulawesi in 2004'),
  (gen_random_uuid(), 'Maluku', 'Moluccas province, historically known as the Spice Islands'),
  (gen_random_uuid(), 'Maluku Utara', 'North Maluku with Ternate and Tidore, centers of ancient spice trade'),
  (gen_random_uuid(), 'Papua Barat', 'West Papua province with rich biodiversity and mining resources'),
  (gen_random_uuid(), 'Papua', 'Papua province with Jayapura as capital, largest province in Indonesia'),
  (gen_random_uuid(), 'Papua Tengah', 'Central Papua, established in 2022 with Nabire as capital'),
  (gen_random_uuid(), 'Papua Pegunungan', 'Highland Papua, established in 2022 with Wamena as capital'),
  (gen_random_uuid(), 'Papua Selatan', 'South Papua, established in 2022 with Merauke as capital'),
  (gen_random_uuid(), 'Papua Barat Daya', 'Southwest Papua, established in 2022 with Sorong as capital');

-- Verify the insertion
SELECT COUNT(*) as total_provinces FROM public.provinces;