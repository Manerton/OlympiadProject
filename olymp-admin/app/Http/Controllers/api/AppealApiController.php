<?php

namespace App\Http\Controllers\api;

use App\Http\Requests\AppealRequest;
use App\Repositories\AppealRepository;
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
        return response()->json([]);
    }
    public function changeStatus(AppealRequest $request, $id){
        $data = $request->validated();
        $this->appealService->changeStatus($id, $data);
        return response()->json([]);
    }
    public function appealByEvent($id)
    {
        $appeals = $this->appealService->getByEventId($id);
        return response()->json([
            'data' => $appeals
        ]);
    }
    public function appealByUser($id){
        $appeals = $this->appealService->getByUserId($id);
        return response()->json([
            'data' => $appeals
        ]);
    }
}
