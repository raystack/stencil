module Stencil
  RSpec.describe Client do
    let(:registry_url) { 'http://stencil.test' }
    let(:bearer_token) { 'sample-token-123' }
    let(:service_type_message)  do
      service_type = nil
      Google::Protobuf::FileDescriptorSet.decode(File.read('spec/data/desc-proto-bin')).file.each do |file|
        file.message_type.each do |msg|
          if msg.name == "ServiceType"
            service_type = msg
          end
        end
      end
      service_type
    end

    context '#get_type' do
      subject { Stencil::Client.new }

      before(:each) do
        Stencil.configure do |config|
          config.registry_url = registry_url
          config.bearer_token = bearer_token
        end

        @stencil_get_stub = stub_request(:get, registry_url).
          with(
            headers: {
              'Authorization' => 'Bearer ' + bearer_token,
              'Connection' => 'close',
              'Host' => 'stencil.test',
              'User-Agent' => 'http.rb/4.4.1'
            })
      end

      it 'should raise error if configs are invalid' do
        config = Stencil.configuration
        config.bearer_token = ""
        expect { subject }.to raise_error(Stencil::InvalidConfiguration)
      end

      it 'should raise error if http client returns error on stencil get api' do
        @stencil_get_stub.to_raise(StandardError.new('some error'))
        expect { subject.get_type }.to raise_error(Stencil::HTTPClientError)
      end

      it 'should raise error if http client returns 500' do
        @stencil_get_stub.to_return(status: 500, body: 'Internal server error', headers: {})
        expect { subject.get_type }.to raise_error(Stencil::HTTPClientError)
      end

      it 'should raise error for invalid proto type' do
        @stencil_get_stub.to_return(status: 200, body: File.new('spec/data/desc-proto-bin'), headers: {})

        proto_name = "incorrect"
        expect { subject.get_type(proto_name) }.to raise_error(Stencil::InvalidProtoClass)
      end

      it 'should successfully return proto type' do
        @stencil_get_stub.to_return(status: 200, body: File.new('spec/data/desc-proto-bin'), headers: {})

        proto_name = "com.gojek.esb.types.ServiceType"
        actual_type = subject.get_type(proto_name)
        expect(actual_type).to eq(service_type_message)
      end
    end
  end
end