package http

import (
	"strconv"

	"poprako-main-server/internal/state"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

func Run(appState state.AppState) {
	// A default Iris application,
	// with logger and recovery middleware already attached.
	app := iris.Default()

	app.Use(logger.New())

	// Apply handlers to routes.
	routeApp(app, &appState)

	// Run the application.
	runServer(app, appState.Cfg.Host, appState.Cfg.Port)
}

// Route application endpoints.
func routeApp(app *iris.Application, appState *state.AppState) {
	api := app.Party("/api/v1")

	// Public routes (no auth required)
	api.Get("/check-update", CheckUpdateHandler(appState))
	api.Post("/login", LoginUser(appState))

	// Apply auth middleware to all routes below
	api.Use(AuthMiddleware(appState))

	users := api.Party("/users")
	{
		users.Get("", RetrieveUserInfos(appState))
		users.Get("/me", GetCurrUserInfo(appState))
		users.Get("/{user_id:string}", GetUserInfoByID(appState))
		users.Get("/invitations", GetInvitations(appState))
		users.Post("/invitations", InviteUser(appState))
		users.Patch("/{user_id:string}/roles", AssignUserRole(appState))
		users.Patch("/{user_id:string}", UpdateUserInfo(appState))
	}

	worksets := api.Party("/worksets")
	{
		worksets.Get("", RetrieveWorksets(appState))
		worksets.Get("/{workset_id:string}", GetWorksetByID(appState))
		worksets.Post("", CreateWorkset(appState))
		worksets.Patch("/{workset_id:string}", UpdateWorksetByID(appState))
		worksets.Delete("/{workset_id:string}", DeleteWorksetByID(appState))
	}

	comics := api.Party("/comics")
	{
		comics.Get("", RetrieveComicBriefs(appState))
		comics.Get("/{comic_id:string}", GetComicInfoByID(appState))
		comics.Get("/{comic_id:string}/export", ExportComic(appState))
		comics.Get("/{comic_id:string}/cover", GetCoverByComicID(appState))
		comics.Get("/{comic_id:string}/pages", GetPagesByComicID(appState))
		comics.Post("/{comic_id:string}/import", ImportComic(appState))
		comics.Post("", CreateComic(appState))
		comics.Patch("/{comic_id:string}", UpdateComicByID(appState))
		comics.Delete("/{comic_id:string}", DeleteComicByID(appState))

		app.HandleDir(appState.ComicSvc.ExportBaseURI(), appState.Cfg.ComicExportDir)
	}

	worksetComics := api.Party("/worksets/{workset_id:string}/comics")
	{
		worksetComics.Get("", GetComicBriefsByWorksetID(appState))
	}

	pages := api.Party("/pages")
	{
		pages.Get("/{page_id:string}", GetPageByID(appState))
		pages.Post("", CreatePages(appState))
		pages.Post("/recreate", RecreatePage(appState))
		pages.Delete("/{page_id:string}", DeletePageByID(appState))
		pages.Patch("/{page_id:string}", UpdatePageByID(appState))
	}

	units := api.Party("/pages/{page_id:string}/units")
	{
		units.Get("", GetUnitsByPageID(appState))
		units.Post("", CreateUnits(appState))
		units.Patch("", UpdateUnits(appState))
		units.Delete("", DeleteUnits(appState))
	}

	asgns := api.Party("/assignments")
	{
		asgns.Get("/{asgn_id:string}", GetAsgnByID(appState))
		asgns.Post("", CreateAsgn(appState))
		asgns.Delete("/{asgn_id:string}", DeleteAsgnByID(appState))
		asgns.Patch("/{asgn_id:string}", UpdateAsgn(appState))
	}

	comicAsgns := api.Party("/comics/{comic_id:string}/assignments")
	{
		comicAsgns.Get("", GetAsgnsByComicID(appState))
	}

	userAsgns := api.Party("/users/{user_id:string}/assignments")
	{
		userAsgns.Get("", GetAsgnsByUserID(appState))
	}
}

func runServer(
	app *iris.Application,
	host string,
	port uint16,
) {
	addr := host + ":" + strconv.Itoa(int(port))

	app.Listen(addr)
}
