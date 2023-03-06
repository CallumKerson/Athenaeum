/**
 NetworkController.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Combine
import Foundation

protocol NetworkController: AnyObject {
    typealias Headers = [String: Any]

    func get<T>(type: T.Type,
                url: URL?,
                headers: Headers) -> AnyPublisher<T, AthenaeumError> where T: Decodable

    func patch<T, V>(
        updateObject: T,
        returnType: V.Type,
        url: URL?,
        headers: Headers
    ) -> AnyPublisher<V, AthenaeumError> where T: Encodable, V: Decodable
}

final class FoundationNetworkController: NetworkController {
    func get<T: Decodable>(type: T.Type, url: URL?,
                           headers: Headers) -> AnyPublisher<T, AthenaeumError>
    {
        guard let url = url else {
            return Fail(
                outputType: type,
                failure: AthenaeumError.message("The URL was not specificed")
            ).eraseToAnyPublisher()
        }

        var urlRequest = URLRequest(url: url)

        headers.forEach { key, value in
            if let value = value as? String {
                urlRequest.setValue(value, forHTTPHeaderField: key)
            }
        }
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        return URLSession.shared.dataTaskPublisher(for: urlRequest)
            .mapError { AthenaeumError.map($0) }
            .map(\.data)
            .decode(type: T.self, decoder: decoder)
            .mapError { AthenaeumError.map($0) }
            .eraseToAnyPublisher()
    }

    func patch<T: Encodable, V: Decodable>(updateObject: T, returnType: V.Type, url: URL?,
                                           headers: Headers) -> AnyPublisher<V, AthenaeumError>
    {
        guard let url = url else {
            return Fail(
                outputType: returnType,
                failure: AthenaeumError.message("The URL was not specificed")
            ).eraseToAnyPublisher()
        }

        var urlRequest = URLRequest(url: url)
        headers.forEach { key, value in
            if let value = value as? String {
                urlRequest.setValue(value, forHTTPHeaderField: key)
            }
        }
        urlRequest.addValue("content/json", forHTTPHeaderField: "Content-Type")
        urlRequest.httpMethod = "PATCH"
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        let body = try? encoder.encode(updateObject)

        urlRequest.httpBody = body

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        return URLSession.shared.dataTaskPublisher(for: urlRequest)
            .mapError { AthenaeumError.map($0) }
            .map(\.data)
            .decode(type: V.self, decoder: decoder)
            .mapError {
                logger.error("An error was emitted from the URL publisher: \($0)")
                return AthenaeumError.map($0)
            }
            .eraseToAnyPublisher()
    }
}
