package views

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Top) Render() wecty.HTML {
	return wecty.Tag("body", 
		wecty.Tag("header", 			
			wecty.Class{
				"navbar": true,
			},
			wecty.Tag("section", 
				wecty.Tag("a", 					
					wecty.Attr("href", "#/"),
					wecty.Class{
						"navbar-brand": true,
					},
					wecty.Text("HS-BAND Utility(for ES1)"),
				),
				wecty.Tag("a", 					
					wecty.Attr("href", "../all/dist/#/"),
					wecty.Class{
						"btn": true,
						"btn-link": true,
					},
					wecty.Text("Recorder"),
				),
			),
		),
		wecty.Tag("main", 			
			wecty.Class{
				"container": true,
			},
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Device("),
							wecty.Tag("span", 
								wecty.Text(c.FirmwareRevision),
							),
							wecty.Text(")"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("button", 								
								wecty.Attr("type", "button"),
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
								},
								wecty.Event("click", c.OnConnect),
								wecty.Text("Connect"),
							),
							wecty.Tag("button", 								
								wecty.Attr("type", "button"),
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Event("click", c.OnDisconnect),
								wecty.Text("Disconnect"),
							),
						),
					),
				),
			),
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnSetLED),
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Attr("for", "color"),
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Set LED Color"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							
							wecty.Tag("input", 								
								wecty.Attr("type", "color"),
								wecty.Attr("name", "color"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("style", "max-width: 400px;"),
							),
							wecty.Tag("button", 								
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Text("SEND"),
							),
						),
					),
				),
			),
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnSetCoreTemp),
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Attr("for", "color"),
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Set Core Temp"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							
							wecty.Tag("input", 								
								wecty.Attr("type", "number"),
								wecty.Attr("name", "coreTemp"),
								wecty.Attr("value", "37.0"),
								wecty.Attr("step", "0.1"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("style", "max-width: 400px;"),
							),
							wecty.Tag("button", 								
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Text("SEND"),
							),
						),
					),
				),
			),
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnShutdown),
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Attr("for", "enterOTA"),
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Enter Shutdown Mode"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("button", 								
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Text("Enter"),
							),
						),
					),
				),
			),
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnFactoryReset),
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Attr("for", "enterOTA"),
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Factory Reset"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("button", 								
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Text("Reset"),
							),
						),
					),
				),
			),
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnEnterOTA),
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Attr("for", "enterOTA"),
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Enter OTA Mode"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("button", 								
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Text("Enter"),
							),
						),
					),
				),
			),
			wecty.Tag("form", 				
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnEnterUF2),
				wecty.Tag("div", 					
					wecty.Class{
						"form-group": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"col-2": true,
							"col-sm-12": true,
						},
						wecty.Tag("label", 							
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Enter UF2 Mode"),
							
							wecty.Tag("br", ),
							wecty.Text("(for v0.3.9 or later)"),
						),
					),
					wecty.Tag("div", 						
						wecty.Class{
							"col-10": true,
							"col-sm-12": true,
						},
						wecty.Tag("div", 							
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("button", 								
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Text("Enter"),
							),
						),
					),
				),
			),
		),
	)
}
