package main

import (
	//"context"
	//"time"

	// "github.com/birukbelay/items/controller/http"
	"fmt"
	"github.com/rs/cors"

	"log"
	"net/http"

	"github.com/joho/godotenv"
	//"github.com/tkanos/gonfig"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"

	//my imports
	"github.com/birukbelay/item/config"
	"github.com/birukbelay/item/config/db"
	"github.com/birukbelay/item/entity"
	mgoProductRepo "github.com/birukbelay/item/models/items/repository/mongo"
	mgoUserRepo "github.com/birukbelay/item/models/user/repository/mongo"
	// userRepo "github.com/birukbelay/items/models/user/repository"

	"github.com/birukbelay/item/controller/http/apiHandler/apiSecurityHandler"
	userService "github.com/birukbelay/item/models/user/service"

	"github.com/birukbelay/item/controller/http/apiHandler/productHandler"
	ZitemService "github.com/birukbelay/item/models/items/services"
)

func createTables(dbconn *gorm.DB) []error {
	errs := dbconn.CreateTable(&entity.Item{}, &entity.Categories{}, &entity.User{}).GetErrors()
	if errs != nil {
		return errs
	}
	return nil
}

func updateTables(dbconn *gorm.DB) []error {
	errs := dbconn.AutoMigrate(&entity.Item{},  &entity.User{}).GetErrors()

	if errs != nil {
		return errs
	}
	return nil
}

func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {

	conf := config.New()


	//config := entity.PostgresConfig{}
	//err := gonfig.GetConf("/config/config.json", &config)
	//p:=fmt.Sprintf("postgres://"+config.PG_USER_NAME+":"+config.PG_PASS+"@localhost/"+config.PG_DB_NAME+"?sslmode=disable")
	//pgDbconn, err := gorm.Open("postgres",p)


	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>|
	//........Postgres database......|
	//...............................|
	// pgDbconn := db.PgDbConn()
	// defer pgDbconn.Close()
	// updateTables(pgDbconn)

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>|
	//......MongoDb database.........|
	//...............................|
	mongoDbConn := db.MongodbConn()
	defer db.CloseMongo(mongoDbConn)
	//........ Collections...........|
	db1 := mongoDbConn.Database("database1")

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>|


	//postgres repos .......................................
	//=========================================================|
	// usrRepo := userRepo.NewUserGormRepo(pgDbconn)
	// sessionRepo := userRepo.NewSessionGormRepo(pgDbconn)
	// roleRepo := userRepo.NewRoleGormRepo(pgDbconn)
	//itemRepo := postgres.NewItemGormRepo(pgDbconn)
	//=========================================================|

	//mongodb repos
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++|

	itemCollection:=db1.Collection("items")
	itemRepo := mgoProductRepo.NewProductMongoRepo(itemCollection)

	categoriesCollection:=db1.Collection("Categories")
	categoriesRepo := mgoProductRepo.NewCategoriesMongoRepo(categoriesCollection)

	userCollection :=db1.Collection("users")
	userRepo := mgoUserRepo.NewUserMongoRepo(userCollection)




	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++|



	userServ := userService.NewUserService(userRepo)


	// roleService := userService.NewRoleService(roleRepo)
	signkey := conf.SignKey
	sh:= apiSecurityhandler.NewUserHandler(userServ,  []byte(signkey))




	itemService := ZitemService.NewItemService(itemRepo)
	itemHandler := productHandler.NewAdminItemHandler(itemService)

	categoriesService := ZitemService.NewCategoriesService(categoriesRepo)
	categoriesHandler := productHandler.NewAdminCategoriesHandler(categoriesService)


	/*===================n  Gin=================*/
	//route := gin.Default()

	//// Serve frontend static files
	//route.Use(static.Serve("/assets", static.LocalFile("./public/assets", true)))
	//
	//// Setup route group for the API
	//api := route.Group("/api")
	//{
	//	api.GET("/", func(c *gin.Context) {
	//		c.JSON(http.StatusOK, gin.H {
	//			"message": "pong",
	//		})
	//	})
	//}
	//
	//items := api.Group("/items")
	//{
	//	items.GET("/", JokeHandler)
	//}

	//========== Gin===============






	router := httprouter.New()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})

	handler := c.Handler(router)
	fmt.Println(handler)

	router.POST("/login", sh.ApiLogin)
	router.POST("/signup", sh.ApiSignup)
	//user
	router.POST("/api/user/create",sh.AdminUsersNew)


	//Products
	router.GET("/api/items", itemHandler.GetItems)
	router.GET("/api/items/filter", itemHandler.GetFilteredItems)
	router.GET("/api/items/get/:id",  itemHandler.GetSingleItem)
	//router.POST("/api/items/create",sh.Authenticated(sh.Authorized(itemHandler.InitiateItem)))
	router.POST("/api/items",itemHandler.CreateItem)
	router.PUT("/api/items/update/:id", itemHandler.UpdateItem)
	router.DELETE("/api/items/delete/:id", itemHandler.DeleteItem)

	//Categoriess
	router.GET("/api/categories", categoriesHandler.GetCategories)
	router.GET("/api/categories/:id",  categoriesHandler.GetSingleCategories)
	//router.POST("/api/items/create",sh.Authenticated(sh.Authorized(itemHandler.InitiateItem)))
	router.POST("/api/categories",categoriesHandler.CreateCategories)
	router.PUT("/api/categories/update/:id", categoriesHandler.UpdateCategories)
	router.DELETE("/api/categories/delete/:id", categoriesHandler.DeleteCategories)



	//api.GET("/jokes", JokeHandler)

	//route.Run(":8080")


	router.ServeFiles("/assets/*filepath", http.Dir("./public/assets"))

	fs := http.FileServer(http.Dir("./public/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	//http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8181", handler)
	fmt.Println("...8181...")



}

func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H {
		"message":"Jokes handler not implemented yet",
	})
}











func startHTTPServer(r http.Handler) *http.Server {
	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("--------- %s", err)
		} else {
			log.Printf("Httpserver: ListenAndServe() closing...")
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}