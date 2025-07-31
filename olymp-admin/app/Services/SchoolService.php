<?php

namespace App\Services;

use App\Builder\SchoolBuilder;
use App\Repositories\SchoolRepository;

class SchoolService
{
    private SchoolRepository $schoolRepository;
    private SchoolBuilder $schoolBuilder;
    public function __construct(
        SchoolRepository $schoolRepository,
        SchoolBuilder $schoolBuilder
    )
    {
        $this->schoolRepository = $schoolRepository;
        $this->schoolBuilder = $schoolBuilder;
    }
    public function find($id){
        $school = $this->schoolBuilder->build($this->schoolRepository->getByApiId($id));
        return $school;
    }
    public function findAll($page = NULL)
    {
        $data = $this->schoolRepository->getByApiAll($page);
        $schools = [];
        foreach ($data as $school) {
            $schools[] = $this->schoolBuilder->build($school);
        }
        return $schools;
    }
}
