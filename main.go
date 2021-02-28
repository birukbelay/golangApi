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

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"

	//my imports
	"github.com/birukbelay/item/config"
	"github.com/birukbelay/item/config/db"
	"github.com/birukbelay/item/entity"
	mgoProductRepo "github.com/birukbelay/item/packages/items/repository/mongo"
	mgoUserRepo "github.com/birukbelay/item/packages/user/repository/mongo"
	// userRepo "github.com/birukbelay/items/packages/user/repository"

	"github.com/birukbelay/item/controller/http/apiHandler/apiSecurityHandler"
	userService "github.com/birukbelay/item/packages/user/service"

	"github.com/birukbelay/item/controller/http/apiHandler/productHandler"
	ZitemService "github.com/birukbelay/item/packages/items/services"
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
	router.POST("/api/user",sh.AdminUsersNew)
	router.PUT("/api/user",sh.AdminUsersUpdate)
	router.DELETE("/api/user",sh.AdminUsersDelete)
	router.GET("/api/users",sh.AdminUsers)
	router.GET("/api/user/:id",sh.User)


	//Products
	router.GET("/api/items", itemHandler.GetItems)
	router.GET("/api/items/filter", itemHandler.GetFilteredItems)
	router.GET("/api/item/:id",  itemHandler.GetSingleItem)
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




	router.ServeFiles("/assets/*filepath", http.Dir("./public/assets"))

	//fs := http.FileServer(http.Dir("./public/assets"))
	//http.Handle("/assets/", http.StripPrefix("/assets", fs))

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8181", handler)
	fmt.Println("...8181...")



}










