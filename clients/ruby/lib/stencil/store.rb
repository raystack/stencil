module Stencil
  class Store
    def initialize
      @lock = Concurrent::ReadWriteLock.new
      @data = Hash.new
    end

    def write(key, value)
      @lock.with_write_lock do
        @data.store(key, value)
      end
    end

    def read(key)
      @data[key]
    end
  end
end