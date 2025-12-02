# BG Dict

A REST server hosting a Bulgarian language dictionary.

Find words and other associated words along with grammar details and pronounciation breakdowns.

Intended as a supporting service for language apps.

## Deploy

The repo provides a docker compose configuration for quickstart.

Create a `.env` file from the provided `.env.example` file. Change any variables to suit your deployment.

Then to start the system run:

`docker compose up`

The system will take a few minutes to start on the first deployment as it downloads the required data and ingests the mysql database dump.

## Acknowledgements

The data used by this system is the database hosted by https://rechnik.chitanka.info/.

This dataset has had translations added to it. These translations were taken from KBEDic Buglarian to English dictionary. Available on sourceforge at https://sourceforge.net/projects/kbedic/.
