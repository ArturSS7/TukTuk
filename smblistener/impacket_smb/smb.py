import logging

from impacket.examples import logger
import smbserver

logger.init(True)
logging.getLogger().setLevel(logging.DEBUG)
server = smbserver.SimpleSMBServer(listenAddress="0.0.0.0", listenPort=445)
server.setSMB2Support(True)
server.setSMBChallenge('')
server.start()