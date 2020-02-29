/**
 Store.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

public final class Store<StoreState: AppState>: ObservableObject where StoreState: Equatable {
    let stateSubject: CurrentValueSubject<StoreState, Never>

    private var dispatchFunction: DispatchFunction!
    private let reducer: Reducer<StoreState>

    public init(reducer: @escaping Reducer<StoreState>,
                middleware: [Middleware<StoreState>] = [],
                state: StoreState) {
        self.reducer = reducer
        self.stateSubject = CurrentValueSubject<StoreState, Never>.init(state)

        var middleware = middleware
        middleware.append(asyncActionsMiddleware)
        self.dispatchFunction = middleware
            .reversed()
            .reduce(
                { [unowned self] action in
                    self._dispatch(action: action) },
                { dispatchFunction, middleware in
                    let dispatch: (Action) -> Void = { [weak self] in self?.dispatch(action: $0) }
                    let getState = { [weak self] in self?.stateSubject.value }
                    return middleware(dispatch, getState)(dispatchFunction)
                }
            )
    }

    public func dispatch(action: Action) {
        DispatchQueue.main.async {
            self.dispatchFunction(action)
        }
    }

    private func _dispatch(action: Action) {
        let currentState = self.stateSubject.value
        let newState = self.reducer(self.stateSubject.value, action)
        if newState != currentState {
            self.stateSubject.send(self.reducer(self.stateSubject.value, action))
        }
    }
}

public let asyncActionsMiddleware: Middleware<AppState> = { dispatch, getState in
    { next in
        { action in
            if let action = action as? AsyncAction {
                action.execute(state: getState(), dispatch: dispatch)
            }
            return next(action)
        }
    }
}
