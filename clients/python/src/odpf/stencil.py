from .store import Store
from schedule import Scheduler
from google.protobuf.message import Message

class MultiUrlClient:
    def __init__(self, urls:list, interval=3600, auto_refresh=False) -> None:
        self._store = Store()
        self._urls = urls
        self._interval = interval
        self._auto_refresh = auto_refresh
        self._schduler = Scheduler()
        if self._auto_refresh:
            self._schduler.every(self._interval).seconds.do(self.refresh)
        #TODO: check whether scheduler executed immediatelly or not
        self.refresh()        

    def refresh(self):
        for url in self._urls:
            self._store.load(url)

    def get_descriptor(self, name:str) -> Message:
        return self._store.get(name)
    
    def parse(self, name:str, data:str):
        msg = self.get_descriptor(name)
        return msg.ParseFromString(data)
    
class Client(MultiUrlClient):
    def __init__(self, url: str) -> None:
        super().__init__([url])