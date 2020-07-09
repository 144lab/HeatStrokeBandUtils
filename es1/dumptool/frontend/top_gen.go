package frontend

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
					wecty.Text("HS-BAND Data-Dump(for ES1)"),
				),
				wecty.Tag("a", 					
					wecty.Attr("href", "../all/dist/#/"),
					wecty.Class{
						"btn": true,
						"btn-link": true,
					},
					wecty.Text("Recorder"),
				),
				wecty.Tag("a", 					
					wecty.Attr("href", "../dist/#/"),
					wecty.Class{
						"btn": true,
						"btn-link": true,
					},
					wecty.Text("Utility"),
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
							wecty.Text("Device"),
							wecty.Tag("span", 
								wecty.Text(c.FirmwareRevision),
							),
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
							wecty.Text("RecordState"),
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
							&c.State,
							wecty.Tag("button", 								
								wecty.Attr("type", "button"),
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Event("click", c.OnReadRecordState),
								wecty.Text("Read"),
							),
						),
					),
				),
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
							wecty.Text("Request Record"),
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
								wecty.Attr("id", "record-start"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("placeholder", "start-id"),
							),
							
							wecty.Tag("input", 								
								wecty.Attr("type", "number"),
								wecty.Attr("id", "record-length"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("placeholder", "length"),
							),
							wecty.Tag("button", 								
								wecty.Attr("type", "button"),
								wecty.Class{
									"btn": true,
									"input-group-btn": true,
									"disabled": true,
								},
								wecty.Event("click", c.OnRequestRecord),
								wecty.Text("Requests"),
							),
						),
					),
				),
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
							wecty.Text("Notify Records"),
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
							wecty.Tag("textarea", 								
								wecty.Attr("id", "log"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("rows", "12"),
							),
						),
					),
				),
			),
		),
	)
}