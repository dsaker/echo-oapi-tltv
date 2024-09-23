# talkliketv

TalkLikeTV is a language learning application designed to address the shortcomings of other language apps by incorporating authentic native speakers conversing at real speeds and tones.

The unique approach involves learning from the subtitles of TV shows and then watching the episodes at your convenience. By translating from your native language to the spoken subtitles of a TV show, you not only grasp how native speakers communicate but also enhance your ability to understand the dialogue of the show.

This api was built using mockgen, sqlc, and oapi-codegen. I was using the net/http library, but I have decided to continue with the echo library because the nethttp-middleware library does not implement a means of passing data from the JWT token down to the handlers.

You can see the issue [here](https://github.com/oapi-codegen/oapi-codegen/blob/b7b82be741ef532eb3f100fb61f62ca3da196ab9/examples/authenticated-api/stdhttp/server/jwt_authenticator.go#L81). I have implemented one way of dealing with this issue [here](https://github.com/dsaker/nethttp-middleware/blob/7a4f0aadf469ca9a655576f46c15b68052a153dc/oapi_validate.go#L113), but have decided to continue with the more mature implementation of [echo-middleware](https://github.com/oapi-codegen/echo-middleware/blob/396ed0328059a05e8a6a47d7af9abb64c735de31/oapi_validate.go#L139)

