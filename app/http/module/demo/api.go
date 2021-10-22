package demo

import (
    demoService "github.com/gohade/hade/app/provider/demo"
    "github.com/gohade/hade/framework/contract"
    "github.com/gohade/hade/framework/gin"
    "github.com/gohade/hade/framework/provider/orm"
)

type DemoApi struct {
    service *Service
}

func Register(r *gin.Engine) error {
    api := NewDemoApi()
    r.Bind(&demoService.DemoProvider{})

    r.GET("/demo/demo", api.Demo)
    r.GET("/demo/demo2", api.Demo2)
    r.POST("/demo/demo_post", api.DemoPost)
    r.GET("/demo/orm", api.DemoOrm)
    return nil
}

func NewDemoApi() *DemoApi {
    service := NewService()
    return &DemoApi{service: service}
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
    c.JSON(200, "this is demo for dev all")
}

// Demo2  for godoc
// @Summary 获取所有学生
// @Description 获取所有学生,不进行分页
// @Produce  json
// @Tags demo
// @Success 200 {array} UserDTO
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
    demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
    students := demoProvider.GetAllStudent()
    usersDTO := StudentsToUserDTOs(students)
    c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoOrm(c *gin.Context) {
    logger := c.MustMakeLog()
    gormService := c.MustMake(contract.ORMKey).(contract.ORMService)
    db, err := gormService.GetDB(orm.WithConfigPath("database.default"))
    if err != nil {
        c.AbortWithError(500, err)
        return
    }

    err = db.AutoMigrate(&User{})
    if err != nil {
        c.AbortWithError(500, err)
        return
    }
    logger.Info(c, "migrate ok", nil)

    c.JSON(200, "ok")
}

func (api *DemoApi) DemoPost(c *gin.Context) {
    type Foo struct {
        Name string
    }
    foo := &Foo{}
    err := c.BindJSON(&foo)
    if err != nil {
        c.AbortWithError(500, err)
    }
    c.JSON(200, nil)
}
