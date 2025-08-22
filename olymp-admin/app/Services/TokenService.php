<?php

namespace App\Services;

use App\Repositories\TokenRepository;

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
}
