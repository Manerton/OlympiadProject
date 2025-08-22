<?php

namespace App\Http\Middleware;

use App\Services\LogService;
use App\Services\TokenService;
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
    private TokenService $tokenService;
    public function __construct(
        LogService $logService,
        TokenService $tokenService
    )
    {
        $this->logService = $logService;
        $this->tokenService = $tokenService;
    }

    public function handle(Request $request, Closure $next)
    {
        $this->logService->create(request()->ip(), request()->userAgent(), request()->url());
        $cookieValue = Cookie::get('username');
        if ($cookieValue) {
            $data = json_decode($cookieValue, true);
            $expireTime = $data['expire'];
            if ($expireTime > time()) {
                return $next($request);
            } else {
                $response = $this->tokenService->refresh();
                Cookie::queue(Cookie::forget('username'));
                Cookie::queue(Cookie::forget('refresh_token'));
                if ($response['success']) {
                    $rawCookie = $response['headers']['Set-Cookie'][0];
                    preg_match('/refresh_token=([^;]+)/', $rawCookie, $matches);
                    $refreshToken = $matches[1] ?? null;
                    $token = $response['data']['data']['access_token'];
                    $expireTime = $response['data']['data']['expires_in'];
                    Cookie::queue('username', json_encode([
                        'email' => $data['email'],
                        'token' => $token,
                        'expire' => now()->addSeconds($expireTime)->timestamp
                    ]), $expireTime);
                    Cookie::queue('refresh_token', $refreshToken, $expireTime);
                    return $next($request);
                }
                else {
                    return redirect('/login');
                }
            }
        }
        else {
            Cookie::queue(Cookie::forget('username'));
            Cookie::queue(Cookie::forget('refresh_token'));
            return redirect('/login');
        }
    }
}
