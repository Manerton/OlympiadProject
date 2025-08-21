<?php

namespace App\Http\Middleware;

use App\Services\LogService;
use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class AuthMiddleware
{
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure(\Illuminate\Http\Request): (\Illuminate\Http\Response|\Illuminate\Http\RedirectResponse)  $next
     * @return \Illuminate\Http\Response|\Illuminate\Http\RedirectResponse
     */
    private LogService $logService;
    public function __construct(
        LogService $logService
    )
    {
        $this->logService = $logService;
    }

    public function handle(Request $request, Closure $next)
    {
        $this->logService->create(request()->ip(), request()->userAgent(), request()->url());
        if (Cookie::get('username')){
            return $next($request);
        }
        else {
            return redirect('/login');
        }
    }
}
