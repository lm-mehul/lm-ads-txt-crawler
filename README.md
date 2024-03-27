# ads-txt-crawler

- Steps to run
- startzookeeper
- startkafka
- run main.go file

BundleIds - 

Android  - com.bundle.similar BundleId
Ios - Integer Bundle Id
CTV - Pixalate - aphanumeric Bundle_Id



main()
|
|-- args = parser.parse_args()
|   |-- bundle_parser() or ads_txt_parser('ads')
|       |-- bundle_parser()
|       |   |-- android_bundle_parser()
|       |   |   |-- requests.get(play_store_url, timeout=5)
|       |   |   |-- BeautifulSoup(response.text, 'html.parser')
|       |   |-- ios_bundle_parser()
|       |   |   |-- requests.head(url, allow_redirects=True, timeout=5)
|       |   |   |-- BeautifulSoup(response.text, 'html.parser')
|       |   |-- ctv_bundle_parser()
|       |       |-- requests.post(algolia_url, headers=headers, data=payload)
|       |-- ads_txt_parser('ads')
|           |-- isAdsTxtLinePresent(ads_txt_page, ads_txt_lines)
|               |-- requests.get(ads_txt_page, headers=headers, timeout=5)
|   |-- main() (if invalid script_type provided)
|
|-- bundle_parser()
|   |-- android_bundle_parser()
|   |-- ios_bundle_parser()
|   |-- ctv_bundle_parser()
|   |-- ads_txt_parser('ads')
|       |-- isAdsTxtLinePresent(ads_txt_page, ads_txt_lines)
|           |-- requests.get(ads_txt_page, headers=headers, timeout=5)
|
|-- ads_txt_parser('ads')
|   |-- isAdsTxtLinePresent(ads_txt_page, ads_txt_lines)
|       |-- requests.get(ads_txt_page, headers=headers, timeout=5)


