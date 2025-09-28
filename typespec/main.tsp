import "@typespec/http";
import "@typespec/rest";
import "@typespec/json-schema";
import "./config/config.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.JsonSchema;

@service({
  title: "GoWatchIt Configuration API"
})
@server("http://localhost:8080", "Development server")
namespace GoWatchIt;

// Configuration models are now directly in GoWatchIt namespace