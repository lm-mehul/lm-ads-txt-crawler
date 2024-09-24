# ads-txt-crawler

- Steps to run
- startzookeeper
- startkafka
- run main.go file

BundleIds - 

Android  - com.bundle.similar BundleId
Ios - Integer Bundle Id
CTV - Pixalate - aphanumeric Bundle_Id

0 17 * * 6 cd /var/data_mount/da_team/lm-ads-txt-crawler && nohup go run main.go --config=resource/config/config.yml >> /var/data_mount/da_team/lm-ads-txt-crawler/logfile.log 2>&1 &


- ![Design Diagram](/resources/Design%20Diagrams/lm-ads-txt-crawler.jpg)
