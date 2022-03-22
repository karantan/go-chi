# github.com/go-chi/chi/v5

Welcome to the gochi generated docs.

## Routes

<details>
<summary>`/`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/**
	- _GET_
		- [main.Hello]()

</details>
<details>
<summary>`/api/v1/articles`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/api/v1**
	- [o-chi/oauth.(*BearerAuthentication).Authorize-fm]()
	- **/articles**
		- **/**
			- _GET_
				- [main.ListArticles]()
			- _POST_
				- [main.CreateArticle]()

</details>
<details>
<summary>`/api/v1/articles/search`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/api/v1**
	- [o-chi/oauth.(*BearerAuthentication).Authorize-fm]()
	- **/articles**
		- **/search**
			- _GET_
				- [main.SearchArticles]()

</details>
<details>
<summary>`/api/v1/articles/{articleID}`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/api/v1**
	- [o-chi/oauth.(*BearerAuthentication).Authorize-fm]()
	- **/articles**
		- **/{articleID}**
			- [main.ArticleCtx]()
			- **/**
				- _DELETE_
					- [main.DeleteArticle]()
				- _GET_
					- [main.GetArticle]()
				- _PUT_
					- [main.UpdateArticle]()

</details>
<details>
<summary>`/api/v1/articles/{articleSlug:[a-z-]+}`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/api/v1**
	- [o-chi/oauth.(*BearerAuthentication).Authorize-fm]()
	- **/articles**
		- **/{articleSlug:[a-z-]+}**
			- _GET_
				- [main.ArticleCtx]()
				- [main.GetArticle]()

</details>
<details>
<summary>`/auth`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/auth**
	- _POST_
		- [o-chi/oauth.(*BearerServer).ClientCredentials-fm]()

</details>
<details>
<summary>`/token`</summary>

- [RequestID]()
- [Logger]()
- [Recoverer]()
- **/token**
	- _POST_
		- [o-chi/oauth.(*BearerServer).UserCredentials-fm]()

</details>

Total # of routes: 7

