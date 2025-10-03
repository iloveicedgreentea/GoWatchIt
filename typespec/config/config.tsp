import "@typespec/http";
import "@typespec/json-schema";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.JsonSchema;
using TypeSpec.Rest;

namespace GoWatchIt;

// Note: Custom decorators would need to be implemented via emitters
// For now, we'll use standard decorators and comments to indicate UI metadata

@doc("EzBEQ integration configuration")
@jsonSchema
model EZBEQConfig {
  @doc("Database ID (internal use only)")
  @visibility("read")
  id?: int64;

  @doc("Enable EzBEQ integration")
  enabled: boolean = false;

  @doc("EzBEQ server hostname or IP address")
  @example("ezbeq.local")
  @pattern("^[a-z0-9.-]+$|^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$")
  url: string = "ezbeq.local";

  @doc("EzBEQ server port")
  @example("8080")
  port: string = "8080";

  @doc("Connection protocol")
  scheme: "http" | "https" = "http";

  @doc("AVR brand")
  @example("Denon")
  avrBrand: "Denon" | "Marantz" = "Denon";

  @doc("AVR URL with protocol")
  @example("http://192.168.1.100")
  @pattern("^https?://")
  avrURL: string = "";

  @doc("Denon AVR IP address")
  @example("192.168.1.100")
  @pattern("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$|^[a-z0-9.-]+$")
  denonIP: string = "";

  @doc("Denon AVR port")
  @example("23")
  denonPort: string = "23";

  @doc("Preferred BEQ authors. Select authors to allow. Leave empty to match any author")
  preferredAuthors: string[] = #[];

  @doc("Profile slot numbers to load into")
  @minItems(0)
  @maxItems(4)
  slots: int32[] = #[];

  @doc("Use AVR codec search for BEQ profiles")
  useAVRCodecSearch: boolean = false;

  @doc("Adjust master volume with profile")
  @jsonSchema
  adjustMasterVolumeWithProfile: boolean = true;

  @doc("Enable TV BEQ profiles")
  enableTVBEQ: boolean = false;

  @doc("Match ANY edition IF your client sends a blank edition")
  looseEditionMatching: boolean = false;

  @doc("Skip edition matching entirely")
  skipEditionMatching: boolean = false;

  @doc("Stop Plex playback if BEQ profile doesn't match")
  stopPlexIfMismatch: boolean = false;

  @doc("Send notification when profile is loaded")
  notifyOnLoad: boolean = false;

  @doc("Send notification when profile is unloaded")
  notifyOnUnload: boolean = false;

  @doc("Enable dry run mode (no actual changes)")
  dryRun: boolean = false;
}

@doc("Home Assistant integration configuration")
@jsonSchema
model HomeAssistantConfig {
  @doc("Database ID (internal use only)")
  @visibility("read")
  id?: int64;

  @doc("Enable Home Assistant integration")
  enabled: boolean = false;

  @doc("Home Assistant URL")
  @example("homeassistant.local")
  @pattern("^[a-z0-9.-]+$|^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$")
  url: string = "";

  @doc("Home Assistant port")
  @example("8123")
  port: string = "8123";

  @doc("Connection protocol")
  scheme: "http" | "https" = "http";

  @doc("Home Assistant authentication token")
  @format("password")
  token: string = "";

  @doc("Remote entity name in Home Assistant")
  @example("remote.shield_tv")
  remoteEntityName: string = "";

  @doc("Notification endpoint name")
  @example("notify.tv")
  notifyEndpointName: string = "";

  @doc("Notification display time in milliseconds")
  @minValue(1000)
  @maxValue(30000)
  notifyDisplayTime: int32 = 5000;

  @doc("Trigger aspect ratio change on media events")
  triggerAspectRatioChangeOnEvent: boolean = false;
}

@doc("Supported media player types")
enum Player {
  Plex: "Plex",
  Jellyfin: "Jellyfin",
  HomeAssistant: "HomeAssistant",
}

@doc("Media player configuration")
@jsonSchema
model PlayerConfig {
  @doc("Database ID (internal use only)")
  @visibility("read")
  id?: int64;

  @doc("Enable player integration")
  enabled: boolean = false;

  @doc("Media player URL")
  @example("plex.local")
  @pattern("^[a-z0-9.-]+$|^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$")
  url: string = "";

  @doc("Media player port")
  @example("32400")
  port: string = "32400";

  @doc("Connection protocol")
  scheme: "http" | "https" = "http";

  @doc("Type of media player")
  playerType: Player = Player.Plex;

  @doc("Authentication token for the media player")
  @format("password")
  token: string = "";

  @doc("User ID for players that require it")
  @example("123456789")
  userID: string = "";

  @doc("Filter by device UUID (comma-separated for multiple)")
  @example("abc123,def456")
  deviceUUIDFilter: string = "";

  @doc("Filter by owner name")
  @example("John Doe")
  ownerNameFilter: string = "";

  @doc("Skip TMDB metadata lookup (Jellyfin only)")
  skipTMDB: boolean = false;

  @doc("Media entity name (Home Assistant only)")
  @example("media_player.living_room")
  mediaEntityName: string = "";
}

@doc("HDMI sync configuration")
@jsonSchema
model HDMISyncConfig {
  @doc("Database ID (internal use only)")
  @visibility("read")
  id?: int64;

  @doc("Enable HDMI sync")
  enabled: boolean = false;

  @doc("HDMI source input")
  @example("HDMI1")
  source: string = "";

  @doc("Sync timing in milliseconds")
  @example("500")
  time: string = "500";

  @doc("Envy device identifier")
  envy: string = "";

  @doc("Media player IP address")
  @example("192.168.1.50")
  @pattern("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$")
  playerIP: string = "";

  @doc("Media player machine identifier")
  @example("abc123def456")
  playerMachineIdentifier: string = "";

  @doc("Connection protocol")
  scheme: "http" | "https" = "http";
}

@doc("BEQ Author from catalogue")
@jsonSchema
model BEQAuthor {
  @doc("Database ID (internal use only)")
  @visibility("read")
  id?: int64;

  @doc("Author name from BEQ catalogue")
  name: string = "";

  @doc("Last updated timestamp (ISO 8601)")
  updatedAt?: string;
}

@doc("Complete application configuration")
@jsonSchema
model AppConfig {
  @doc("EzBEQ integration settings")
  ezbeq: EZBEQConfig;

  @doc("Home Assistant integration settings")
  homeAssistant: HomeAssistantConfig;

  @doc("Media player settings")
  player: PlayerConfig;

  @doc("HDMI sync settings")
  hdmiSync: HDMISyncConfig;
}

// REST API endpoints to expose models in OpenAPI
@route("/config")
interface ConfigAPI {
  @get
  @doc("Get the complete application configuration")
  getConfig(): AppConfig;

  @put
  @doc("Update the complete application configuration")
  updateConfig(@body config: AppConfig): AppConfig;

  @route("/ezbeq")
  @get
  @doc("Get EzBEQ configuration")
  getEzbeqConfig(): EZBEQConfig;

  @route("/homeassistant")
  @get
  @doc("Get Home Assistant configuration")
  getHaConfig(): HomeAssistantConfig;

  @route("/player")
  @get
  @doc("Get player configuration")
  getPlayerConfig(): PlayerConfig;

  @route("/hdmi-sync")
  @get
  @doc("Get HDMI sync configuration")
  getHdmiSyncConfig(): HDMISyncConfig;

  @route("/authors")
  @get
  @doc("Get available BEQ authors from catalogue")
  getAuthors(): string[];
}
