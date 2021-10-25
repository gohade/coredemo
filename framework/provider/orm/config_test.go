package orm

//
//func TestHadeConfig_Load(t *testing.T) {
//	container := tests.InitBaseContainer()
//	container.Bind(&config.HadeConfigProvider{})
//
//	Convey("test config", t, func() {
//		configService := container.MustMake(contract.ConfigKey).(contract.Config)
//		config := &Config{}
//		err := configService.Load("database.default", config)
//		So(err, ShouldBeNil)
//	})
//
//	Convey("test default config", t, func() {
//		configService := container.MustMake(contract.ConfigKey).(contract.Config)
//		config := &Config{
//			ConnMaxIdle: 10,
//		}
//		err := configService.Load("database.read", config)
//		So(err, ShouldBeNil)
//		So(config.ConnMaxIdle, ShouldEqual, 10)
//	})
//
//	Convey("test base config", t, func() {
//		configService := container.MustMake(contract.ConfigKey).(contract.Config)
//		config := &Config{
//			ConnMaxOpen: 200,
//		}
//		err := configService.Load("database", config)
//		So(err, ShouldBeNil)
//		So(config.ConnMaxOpen, ShouldEqual, 100)
//	})
//
//}
