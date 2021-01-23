package external

/*
curl --request POST \
  --url http://bigd.oindata.cn/data-api/api/bookSet \
  --header 'Content-Type: application/json' \
  --header 'appKey: mZ0z1Df0BirPWnrvGVBg5FKqR_B-uDVOnIlTQbofexQ' \
  --data '{ \
  "scenic_id":"8195d516a9374a42bdf88d129eeb13ca", \
 "max_capacity":10000,
 "realtime_capacity":2000,
  "booktime_set":[
      {
       "label":"08:00-10:00",
       "maxBook":2000,
       "starkClock":"08:00",
       "endClock":"10:00"
      },
      {
       "label":"10:00-12:00",
       "maxBook":2000,
       "starkClock":"10:00",
       "endClock":"12:00"
      },
      {
       "label":"12:00-14:00",
       "maxBook":2000,
       "starkClock":"12:00",
       "endClock":"14:00"
      },
      {
       "label":"14:00-16:00",
       "maxBook":2000,
       "starkClock":"14:00",
       "endClock":"16:00"
      }
  ]
}
'*/
