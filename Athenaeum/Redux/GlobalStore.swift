/**
 Impl.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

let store = Store<GlobalAppState>(reducer: appStateReducer, middleware: [logMiddleware],
                                  state: GlobalAppState())
