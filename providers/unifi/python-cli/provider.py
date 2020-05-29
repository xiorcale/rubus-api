from dataclasses import dataclass
from subprocess import PIPE, Popen

from flask import Flask, jsonify

app = Flask(__name__)


@dataclass
class Device:
    '''Device is compliant to the description of a Rubus Device.'''
    id: int
    isTurnedOn: int
    hostname: str


@app.route('/device', methods=['GET'])
def get_all_devices():
    '''Return the list of devices available for Rubus to provision.'''
    with Popen(['./poe.sh', '-c'], stdout=PIPE) as p:
        devices = []
        for line in p.stdout:
            device_raw = str(line, 'utf-8').split(' ')
            device = Device(
                int(device_raw[0]),
                True if device_raw[1] == 'Up' else False,
                device_raw[2].strip()
            )
            devices.append(device)

    return jsonify(devices)


@app.route('/device/<int:id>', methods=['GET'])
def get_device(id: int):
    '''Return a single device information.'''
    with Popen(['./poe.sh', '-c', '-p', str(id)], stdout=PIPE) as p:
        output, _ = p.communicate()
        device_raw = str(output, 'utf-8').split(' ')

    device = Device(
        int(device_raw[0]),
        True if device_raw[1] == 'Up' else False,
        device_raw[2].strip()
    )

    return jsonify(device)


@app.route('/device/<int:id>/on', methods=['POST'])
def power_on(id: int):
    '''Turn on the given device.'''
    with Popen(['./poe.sh', '-u', '-p', str(id)], stdout=PIPE) as p:
        p.communicate()

    return '', 204


@app.route('/device/<int:id>/off', methods=['POST'])
def power_off(id: int):
    '''Turn off the given device.'''
    with Popen(['./poe.sh', '-d', '-p', str(id)], stdout=PIPE) as p:
        p.communicate()

    return '', 204


if __name__ == "__main__":
    app.run(host='0.0.0.0', debug=True, port=1080)
