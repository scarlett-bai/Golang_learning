import { Coolcar } from './request';
import { car } from "./proto_gen/car/car_pb";
import camelcaseKeys = require('camelcase-keys');

export namespace CarService {
    export function subscrbe(onMsg: (c: car.v1.ICarEntity) => void) {
        const socket = wx.connectSocket({
            url: Coolcar.wsAddr + '/ws'
        })
        socket.onMessage(msg => {
            const obj = JSON.parse(msg.data as string)
            onMsg(car.v1.CarEntity.fromObject(
                camelcaseKeys(obj, {
                    deep: true,
                })
            ))
        })
        return socket
    }

    export function getCar(id: string): Promise<car.v1.ICar> {
        return Coolcar.sendRquestWithAuthRetry({
            method: 'GET',
            path: `/v1/car/${encodeURIComponent(id)}`,
            respMarshaller: car.v1.Car.fromObject,
        })
    }
}