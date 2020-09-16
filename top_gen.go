package main

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("div", 			
			spago.A("class", spago.S(`card`)),
			spago.Tag("div", 				
				spago.A("class", spago.S(`card-header`)),
				spago.Tag("div", 					
					spago.A("class", spago.S(`card-title h5`)),
					spago.T(`JS2Go`),
				),
			),
			spago.Tag("div", 				
				spago.A("class", spago.S(`card-body`)),
				spago.Tag("form", 					
					spago.Event("submit", c.OnSubmit),
					spago.A("class", spago.S(`columns`)),
					spago.Tag("div", 						
						spago.A("class", spago.S(`column col-6`)),
						spago.Tag("div", 							
							spago.A("class", spago.S(`float-right`)),
							spago.Tag("button", 								
								spago.A("class", spago.S(`btn btn-primary`)),
								spago.T(`Convert`),
								spago.Tag("i", 									
									spago.A("class", spago.S(`icon icon-forward`)),
								),
							),
						),
						spago.Tag("textarea", 							
							spago.A("class", spago.S(`form-input`)),
							spago.A("name", spago.S(`js`)),
							spago.A("rows", spago.S(`12`)),
							spago.T(``, spago.S(c.JsCode), ``),
						),
					),
					spago.Tag("div", 						
						spago.A("class", spago.S(`column col-6`)),
						spago.Tag("div", 							
							spago.A("class", spago.S(`float-right`)),
							spago.Tag("button", 								
								spago.A("class", spago.S(`btn btn-primary`)),
								spago.T(`Copy`),
							),
						),
						spago.Tag("textarea", 							
							spago.A("class", spago.S(`form-input`)),
<<<<<<< HEAD
<<<<<<< HEAD
							spago.A("rows", spago.S(`12`)),
							spago.A("readonly", spago.S(`true`)),
							spago.T(``, spago.S(c.GoCode), ``),
=======
							spago.A("rows", spago.S(`8`)),
=======
							spago.A("rows", spago.S(`12`)),
>>>>>>> 686d81a... wip
							spago.A("readonly", spago.S(`true`)),
<<<<<<< HEAD
>>>>>>> cdc3d9a... es5 support completed
=======
							spago.T(``, spago.S(c.GoCode), ``),
>>>>>>> e15b2fb... improve
						),
					),
				),
			),
		),
	)
}
