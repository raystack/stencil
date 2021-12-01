require_relative "stencil/version"
require_relative "stencil/configuration"
require_relative "stencil/constants"
require_relative "stencil/client"

require "http"
require "concurrent"
require "protobuf"

module Stencil
  class Error < StandardError; end
  class InvalidConfiguration < Error; end
  class InvalidProtoClass < Error; end
  class HTTPClientError < Error; end
  class HTTPServerError < Error; end
end
