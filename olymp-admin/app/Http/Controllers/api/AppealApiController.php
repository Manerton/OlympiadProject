<?php

namespace App\Http\Controllers\api;

use App\Http\Requests\AppealRequest;
use App\Services\AppealService;


class AppealApiController
{
    private AppealService $appealService;
    public function __construct(
        AppealService $appealService
    )
    {
        $this->appealService = $appealService;
    }

    public function store(AppealRequest $request){
        $data = $request->validated();
        $this->appealService->create($data);
    }
    public function changeStatus(AppealRequest $request, $id){
        $data = $request->validated();
        $this->appealService->changeStatus($id, $data);
    }
}
