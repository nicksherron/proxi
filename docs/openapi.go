/*
 * Copyright © 2020 nicksherron <nsherron90@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package docs

var doc = `
{
  "openapi": "3.0.0",
  "servers": [
    {
      "url": "{{.Host}}"
    }
  ],
  "info": {
    "description": "Downloads,  checks and stores proxies from the web with rest api for querying results.",
    "version": "{{.Version}}",
    "title": "ProxyPool",
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    }
  },
  "paths": {
    "/get": {
      "get": {
        "summary": "Return a proxy that passed checks.",
        "description": "",
        "operationId": "getProxyWithLimit",
        "parameters": [
          {
            "name": "country",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "allowEmptyValue": true,
            "description": "Filter by country. Format is 'US', 'CH' etc."
          },
          {
            "name": "anon",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean"
            },
            "description": "Only return proxies that where found to be anonymous from tests.  Only need to be present in query params to be true, eg /get?anon",
            "allowEmptyValue": true
          }
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Proxy"
                }
              }
            }
          }
        }
      }
    },
    "/get/{limit}": {
      "get": {
        "summary": "Return n proxies that passed checks.",
        "description": "",
        "operationId": "getNProxyWithLimit",
        "parameters": [
          {
            "name": "country",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "allowEmptyValue": true,
            "description": "Filter by country. Format is 'US', 'CH' etc."
          },
          {
            "name": "anon",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean"
            },
            "description": "Only return proxies that where found to be anonymous from tests.  Only need to be present in query params to be true, eg /get?anon",
            "allowEmptyValue": true
          },
          {
            "name": "limit",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            },
            "description": "Number of proxies to return."
          }
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProxyArray"
                }
              }
            }
          }
        }
      }
    },
    "/getall": {
      "get": {
        "summary": "Return all proxies ignoring filters or status. Warning! may produce lots of results",
        "description": "",
        "operationId": "getWithNoLimit",
        "parameters": [
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProxyArray"
                }
              }
            }
          }
        }
      }
    },
    "/delete": {
      "post": {
        "summary": "Delete a proxy.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Find"
              }
            }
          }
        },
        "operationId": "",
        "responses": {
          "default": {
            "description": "Default response"
          }
        }
      }
    },
    "/find": {
      "post": {
        "summary": "Find proxy.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Find"
              }
            }
          }
        },
        "operationId": "",
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ProxyArray"
                }
              }
            }
          }
        }
      }
    },
    "/stats": {
      "get": {
        "summary": "Shows stats on db and proxies, including the number found, checked, timed out, good and anonymous.",
        "parameters": [
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Stats"
                }
              }
            }
          }
        }
      }
    },
    "/refresh": {
      "get": {
        "summary": "Re-download and check proxies if the server is not already busying downloading or checking. Returns busy, if so.",
        "parameters": [
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
            }
          }
        }
      }
    },
    "/busy": {
      "get": {
        "summary": "Checks whether server is busy with downloads or checks.",
        "parameters": [
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "content": {
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Stats": {
        "type": "object",
        "properties": {
          "total": {
            "type": "integer",
            "example": 62594
          },
          "timeout": {
            "type": "integer",
            "example": 17360
          },
          "anon": {
            "type": "integer",
            "example": 1727
          },
          "recently-checked": {
            "type": "integer",
            "example":1024
          },
          "good": {
            "type": "integer",
            "example": 1727
          }
        }
      },
      "Proxy": {
        "type": "object",
        "properties": {
          "id": {
            "type": "integer",
            "format": "int64",
            "example": 21039
          },
          "created_at": {
            "type": "string",
            "example":"2020-01-27T15:06:18.872358-05:00"
          },
          "updated_at": {
            "type": "string",
            "example":"2020-01-28T04:57:06.613106-05:00"
          },
          "check_count": {
            "type": "integer",
            "example": 3
          },
          "country": {
            "type": "string",
            "example":"IN"
          },
          "fail_count": {
            "type": "integer",
            "example": 0
          },
          "last_status": {
            "type": "string",
            "example": "good"
          },
          "proxy": {
            "type": "string",
            "example": "http://59.91.121.113:35665"
          },
          "timeout_count": {
            "type": "integer",
            "example": 1
          },
          "source": {
            "type": "string",
            "example": "blogspot.com"
          },
          "success_count": {
            "type": "integer",
            "example": 2
          },
          "anonymous": {
            "type": "boolean",
            "example": true
          }
        }
      },
      "ProxyArray": {
        "type": "array",
        "items": {
          "$ref": "#/components/schemas/Proxy"
        }
      },
      "Find": {
        "type": "object",
        "properties": {
          "proxy": {
            "type": "string"
          }
        }
      }
    },
    "requestBodies": {
      "Proxy": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Proxy"
            }
          }
        }
      }
    },
    "links": {},
    "callbacks": {}
  }
}
`
