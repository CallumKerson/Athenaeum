/**
 Protocols.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

public protocol Action: Codable {}

public protocol AsyncAction: Action {
    func execute(state: AppState?, dispatch: @escaping DispatchFunction)
}

public protocol AppState {}

public typealias DispatchFunction = (Action) -> Void
public typealias Middleware<AppState> = (@escaping DispatchFunction, @escaping () -> AppState?)
    -> (@escaping DispatchFunction) -> DispatchFunction
public typealias Reducer<AppState> = (_ state: AppState, _ action: Action) -> AppState
