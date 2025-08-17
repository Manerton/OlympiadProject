<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Prometheus\CollectorRegistry;
use Prometheus\Exception\MetricsRegistrationException;
use Prometheus\Storage\Redis;

class PrometheusMiddleware
{
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure(\Illuminate\Http\Request): (\Illuminate\Http\Response|\Illuminate\Http\RedirectResponse)  $next
     * @return \Illuminate\Http\Response|\Illuminate\Http\RedirectResponse
     */

    protected CollectorRegistry $registry;

    public function __construct()
    {
        $adapter = new Redis([
            'host' => env('REDIS_HOST', 'redis'),
            'port' => env('REDIS_PORT', 6379),
        ]);

        $this->registry = new CollectorRegistry($adapter);
    }

    public function handle($request, Closure $next)
    {
        $response = $next($request);

        try {
            $counter = $this->registry->registerCounter(
                'app',
                'http_requests_total',
                'Total HTTP requests',
                ['method', 'path', 'status']
            );
        } catch (MetricsRegistrationException $e) {
            $counter = $this->registry->getCounter('app', 'http_requests_total');
        }

        $counter->inc([
            $request->method(),
            $request->path(),
            $response->getStatusCode()
        ]);

        return $response;
    }
}
