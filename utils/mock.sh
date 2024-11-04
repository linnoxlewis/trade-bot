PROJECT_PATH=../internal

REPOSITORY_PATH="${PROJECT_PATH}/service"

if ! command -v mockgen &> /dev/null; then
    echo "mockgen не установлен. Установите его, выполнив: go get github.com/golang/mock/mockgen"
    exit 1
fi

for folder in "${REPOSITORY_PATH}"/*; do
    if [ -d "${folder}" ]; then
        PACKAGE_NAME=$(basename "${folder}")

        MOCKS_DIR="${REPOSITORY_PATH}/${PACKAGE_NAME}/mocks"
        mkdir -p "${MOCKS_DIR}"

        mockgen -destination="${MOCKS_DIR}/mock_${PACKAGE_NAME}.go" -package="mocks" "${PROJECT_PATH}/internal/service/${PACKAGE_NAME}" YourInterfaceToMock
        echo "Создан мок для ${PACKAGE_NAME}"
    fi
done