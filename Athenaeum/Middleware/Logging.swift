/**
 Logging.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

let logMiddleware: Middleware<AppState> = { _, _ in
    { next in
        { action in
            #if DEBUG
                let name = __dispatch_queue_get_label(nil)
                let queueName = String(cString: name, encoding: .utf8)
                if let queueName = queueName {
                    log.verbose("Received action \(action) on queue \(queueName)")
                } else {
                    log.verbose("Received action \(action)")
                }
            #endif
            return next(action)
        }
    }
}
