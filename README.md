# TalkLikeTv

TalkLikeTv is a language learning application designed to address limitations I’ve encountered in other popular language learning apps, such as Pimsleur, Babbel, and Duolingo. While these tools serve as strong foundational resources, I’ve found that they tend to plateau once reaching an intermediate level. Currently, I can understand French and Spanish well enough to follow audiobooks and read at a high level, but I still face challenges in expressing myself and comprehending native speakers during travel.

To overcome these barriers, I’ve created an application that generates a Pimsleur-like audio course from any file the user selects. Personally, I use subtitles from current TV shows from the countries I plan to visit. This approach has several benefits: it familiarizes me with contemporary slang, improves my understanding of spoken dialogue, and challenges me to express myself more naturally. Practicing with these audio files not only enhances comprehension of the shows but also provides an immersive, effective way to advance my language skills.

### Installation

- [Install Docker](https://docs.docker.com/engine/install/)
- [Install GoLang](https://go.dev/doc/install)
- Create [Google Cloud Account](https://console.cloud.google.com/getting-started?pli=1)
- Install the [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- Setup [GCP ADC](https://cloud.google.com/docs/authentication/external/set-up-adc )
- Create a [Google Cloud Project](https://developers.google.com/workspace/guides/create-project)
- Run below commands to sign in and enable necessary Google Cloud API's
- Install [ffmpeg](https://www.ffmpeg.org/download.html)
- Install [migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)

```
gcloud init
gcloud services enable texttospeech.googleapis.com
gcloud services enable translate.googleapis.com
```
- Run below commands to start the application
```
git clone https://github.com/dsaker/echo-oapi-tltv.git 
cd echo-oapi-tltv
go mod tidy
docker pull postgres
docker run -d -P -p 127.0.0.1:5433:5432 -e POSTGRES_PASSWORD="password" --name talkliketvpg postgres
echo "export TLTV_DB_DSN=postgresql://postgres:password@localhost:5433/postgres?sslmode=disable" >> .envrc
make db/migrations/up
make run
```
- open http://localhost:8080/swagger/ in local browser
- click on POST /audio/fromfile and click on "Try it out"
- 