PUT /goods/_doc/1
{
    "subs": [
        {
            "id": 158,
            "name": "观光车",
            "pic": "http://data.xytschool.com/storage/image/20200425/1587783884965818.png",
            "goods_type": "ticket",
            "num": 1
        },
        {
            "id": 171,
            "name": "普通标准单人间",
            "pic": "http://data.xytschool.com/storage/image/20200421/1587435808985623.png",
            "goods_type": "room",
            "num": 1
        }
    ],
    "props": [
        {
            "name": "color",
            "subs": [
                "red",
                "blue",
                "black"
            ]
        },
        {
            "name": "size",
            "subs": [
                "sm",
                "m",
                "l"
            ]
        },
        {
            "name": "weight",
            "subs": [
                "500g",
                "1000g",
                "2000g"
            ]
        }
    ],
    "skus": [
        {
            "id": 10,
            "combine": [
                "red",
                "sm",
                "500g"
            ],
            "price": 180.02
        },
        {
            "id": 10,
            "combine": [
                "red",
                "sm",
                "1000g"
            ],
            "price": 280.02
        }
    ]
}

GET goods/_mapping
DELETE goods
PUT goods
{
  "mappings": {
    "properties": {
      "cates": {
        "type": "nested",
        "properties": {
          "name": {
            "type": "keyword"
          },
          "children": {
            "type": "keyword"
          }
        }
      },
      "skus": {
        "type": "nested",
        "properties": {
          "combine": {
            "type": "keyword"
          },
          "id": {
            "type": "long"
          },
          "price": {
            "type": "float"
          }
        }
      },
      "sub_goods": {
        "type": "nested",
        "properties": {
          "goods_type": {
            "type": "keyword"
          },
          "id": {
            "type": "long"
          },
          "name": {
            "type": "keyword"
          },
          "num": {
            "type": "long"
          },
          "pic": {
            "type": "keyword",
            "index": false
          }
        }
      }
    }
  }
}
