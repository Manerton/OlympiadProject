<?php

namespace App\Services;

use App\Repositories\TokenRepository;
use Illuminate\Support\Facades\Cookie;

class TokenService
{
    private TokenRepository $tokenRepository;
    public function __construct(
        TokenRepository $tokenRepository
    )
    {
        $this->tokenRepository = $tokenRepository;
    }
    public function revoke($userId)
    {
        $this->tokenRepository->revoke($userId);
    }
    public function refresh()
    {
        return $this->tokenRepository->refresh();
    }
    public function getToken()
    {
        $token = null;
        if (isset(Cookie::getQueuedCookies()[0])) {
            $queuedCookie = Cookie::getQueuedCookies()[0];
            $token = json_decode($queuedCookie->getValue(), true)['token'] ?? null;
        }
        if (is_null($token) && Cookie::get('username')) {
            $token = json_decode(Cookie::get('username'), true)['token'] ?? null;
        }
        if(is_null($token)){
            $token = str_replace('Bearer ', '', request()->header('Authorization'));
        }
        return $token;
    }
}
