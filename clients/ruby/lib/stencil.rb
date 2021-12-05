require_relative "stencil/version"
require_relative "stencil/configuration"
require_relative "stencil/constants"
require_relative "stencil/client"
require_relative "stencil/store"

require "http"
require "concurrent/timer_task"
require "concurrent/mutable_struct"
require "protobuf"

module Stencil
  class Error < StandardError; end
  class InvalidConfiguration < Error; end
  class InvalidProtoClass < Error; end
  class HTTPClientError < Error; end
  class HTTPServerError < Error; end
end
