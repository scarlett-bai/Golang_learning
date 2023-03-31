import { rental } from "./proto_gen/rental/rental_pb";
import { Coolcar } from "./request";

export namespace TripService {
    export function CreateTrip(req: rental.v1.ICreateTripRequest): Promise<rental.v1.ITripEntity> {
        return Coolcar.sendRquestWithAuthRetry({
            method: 'POST',
            path: '/v1/trip',
            data: req,
            respMarshaller: rental.v1.TripEntity.fromObject,
        })
    }

    export function GetTrip(id: string): Promise<rental.v1.ITrip> {
        return Coolcar.sendRquestWithAuthRetry({
            method: 'GET',
            path: `/v1/trip/${encodeURIComponent(id)}`,
            respMarshaller: rental.v1.Trip.fromObject,
        })
    }

    export function GetTrips(s?: rental.v1.TripStatus): Promise<rental.v1.IGetTripsRequest> {
        let path = '/v1/trips'
        if (s) {
            path += `?status=${s}`
        }
        return Coolcar.sendRquestWithAuthRetry({
            method: 'GET',
            path: path,
            respMarshaller: rental.v1.GetTripsRequest.fromObject,
        })
    }
}