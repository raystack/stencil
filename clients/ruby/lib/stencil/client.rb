module Stencil
  class Client
    attr_reader :root
    def initialize
      begin
        @config = Stencil.configuration
        validate_configuration(@config)
        @root = nil
        setup_http_client
        load_descriptors
      end
    end

    def get_type(proto_name)
      @root.file.each do |file_desc|
        file_desc.message_type.each do |message|
          if proto_name == "#{file_desc.options.java_package}.#{message.name}"
            return message
          end
        end
      end
      raise InvalidProtoClass.new
    end

    private

    def validate_configuration(configuration)
      raise Stencil::InvalidConfiguration.new() if configuration.registry_url.nil? || configuration.bearer_token.nil? || configuration.bearer_token == "Bearer "
    end

    def setup_http_client
      @http_client = HTTP.auth(@config.bearer_token).timeout(@config.http_timeout)
    end

    def load_descriptors
      begin
        response = @http_client.get(@config.registry_url)
        if response.code != 200
          raise HTTPServerError.new("Error while fetching descriptor file: Got #{response.code} from stencil server")
        end
      rescue StandardError => e
        raise HTTPClientError.new(e.message)
      end

      @root = Google::Protobuf::FileDescriptorSet.decode(response.body)
    end
  end
end