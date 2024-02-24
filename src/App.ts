import log from 'loglevel'
import { FitnessMachine } from './fitnessmachine/FitnessMachine'
import { Monitor } from './monitor/Monitor';
import { Capture } from './monitor/Capture';

log.setLevel('DEBUG')

const fitnessMachine = new FitnessMachine({ name: "FDF Rower" })
const monitor = new Monitor({ port: "/dev/tty.usbserial-A6029KI0" }, (captur: Capture) => {
  fitnessMachine.onCapture(captur);
});
monitor.connect((error?) => {
  if (error) {
    process.exit(1);
  }
  fitnessMachine.start();
});

process.on('SIGINT', function () {
  let exitCode = 1;
  monitor.disconnect((error?) => {
    if (!error) {
      exitCode = 0;
    }
  })
  fitnessMachine.stop();

  setTimeout(() => {
    log.info("Bye Bye");
    process.exit(exitCode);
  }, 3000);
});
