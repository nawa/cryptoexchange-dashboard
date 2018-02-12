package main

import (
	i18n "github.com/iris-contrib/middleware/go-i18n"
	"github.com/kataras/iris"
)

type User struct {
	Name string
}

func main() {
	app := iris.New()

	locale := i18n.New(
		"./locales/en-US.all.yaml",
		"./locales/en-US.errors.yaml",
		"./locales/zh-CN.all.yaml",
	)

	app.Use(locale.Serve)

	app.Get("/hi", func(ctx iris.Context) {

		// it tries to find the language by:
		// ctx.Values().GetString("language")
		// if that was empty then
		// it tries to find from the URLParameter setted on the configuration
		// if not found then
		// it tries to find the language by the "language" cookie
		// if didn't found then it it set to the Default setted on the configuration

		// hi is the key, 'iris' is the %s on the .ini file
		// the second parameter is optional

		// hi := ctx.Translate("hi", context.Map{
		// 	"Name": "Peter",
		// })
		hi := ctx.Translate("hi", &User{
			Name: "Peter",
		})

		language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
		// return is form of 'en-US'
		// fmt.Printf("trans %v\n", locale.LanguageTranslationIDs(language))
		// The first succeed language found saved at the cookie with name ("language"),
		//  you can change that by changing the value of the:  iris.TranslateLanguageContextKey
		ctx.Writef("From the language %s translated output: %s", language, hi)
	})

	// go to http://localhost:8080/?lang=el-GR
	// or http://localhost:8080
	// or http://localhost:8080/?lang=zh-CN
	app.Run(iris.Addr(":8080"), iris.WithoutVersionChecker)

}
