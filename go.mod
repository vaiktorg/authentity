module github.com/vaiktorg/Authentity

go 1.18

replace github.com/vaiktorg/gwt => github.com/Vaiktorg/gwt v0.0.0-20220413023908-ee712fe3164d

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/vaiktorg/grimoire v0.0.0-20220112015009-252115daf0b4
	github.com/vaiktorg/gwt v0.0.0-20220413023908-ee712fe3164d
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88
	gorm.io/driver/sqlite v1.3.2
	gorm.io/gorm v1.23.5
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
