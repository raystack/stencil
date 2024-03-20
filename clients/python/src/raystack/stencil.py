from schedule import Scheduler
from google.protobuf.message import Message

from raystack.store import Store


class MultiUrlClient:
    def __init__(self, urls:list, interval=3600, auto_refresh=False) -> None:
        self._store = Store()
        self._urls = urls
        self._interval = interval
        self._auto_refresh = auto_refresh
        self._scheduler = Scheduler()
        if self._auto_refresh:
            self._scheduler.every(self._interval).seconds.do(self.refresh)
        self.refresh()

    def refresh(self):
        for url in self._urls:
            self._store.load(url=url)

    def get_descriptor(self, name: str) -> Message:
        return self._store.get(name)

    def parse(self, name: str, data: bytes):
        msg = self.get_descriptor(name)
        return msg.ParseFromString(data)


class Client(MultiUrlClient):
    def __init__(self, url: str) -> None:
        super().__init__([url])
