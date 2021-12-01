module Stencil
  class Configuration
    def initialize
      @config = ::OpenStruct.new
    end

    def registry_url
      @config.registry_url
    end

    def registry_url=(registry_url)
      @config.registry_url = registry_url
    end

    def http_timeout
      @config.http_timeout || DEFAULT_TIMEOUT_IN_MS
    end

    def http_timeout=(timeout)
      @config.http_timeout = timeout
    end

    def refresh_enabled
      @config.refresh_enabled.nil? ? true : @config.refresh_enabled
    end

    def refresh_enabled=(refresh_enabled = true)
      @config.refresh_enabled = refresh_enabled
    end

    def refresh_ttl_in_secs
      @config.refresh_ttl_in_secs || DEFAULT_REFRESH_INTERVAL_IN_SECONDS
    end

    def refresh_ttl_in_secs=(refresh_ttl_in_secs)
      @config.refresh_ttl_in_secs = refresh_ttl_in_secs
    end

    def bearer_token=(token)
      @config.bearer_token = "Bearer " + token
    end

    def bearer_token
      @config.bearer_token
    end
  end

  def self.configuration
    @config ||= Configuration.new
  end

  def self.configure
    yield(configuration)
  end
end